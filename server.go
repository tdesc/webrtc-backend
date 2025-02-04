package main

import (
	"context"
	"log"
	"sync"
	"time"

	pb "github.com/tdesc/webrtc-backend/sessionpb"

	"github.com/pion/webrtc/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// sessions stores active PeerConnections indexed by session ID.
var sessions = make(map[string]*webrtc.PeerConnection)
var sessionsLock sync.Mutex

// server implements pb.SessionServiceServer.
type server struct {
	pb.UnimplementedSessionServiceServer
}

// StartSession creates a new WebRTC session with a dummy video track and returns the SDP offer.
func (s *server) StartSession(ctx context.Context, req *pb.SessionRequest) (*pb.SessionResponse, error) {
	// Create a new PeerConnection using default configuration.
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return nil, err
	}

	// Create a dummy video track (in a real application, you would stream actual media).
	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
	if err != nil {
		return nil, err
	}
	if _, err = peerConnection.AddTrack(videoTrack); err != nil {
		return nil, err
	}

	// Create an SDP offer.
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		return nil, err
	}
	if err = peerConnection.SetLocalDescription(offer); err != nil {
		return nil, err
	}

	// Store the session.
	sessionsLock.Lock()
	sessions[req.SessionId] = peerConnection
	sessionsLock.Unlock()

	return &pb.SessionResponse{
		SessionId:   req.SessionId,
		WebrtcOffer: offer.SDP,
	}, nil
}

// ConnectSession allows a client to connect to an existing session.
// For simplicity, this returns the existing SDP offer.
func (s *server) ConnectSession(ctx context.Context, req *pb.SessionRequest) (*pb.SessionResponse, error) {
	sessionsLock.Lock()
	peerConnection, exists := sessions[req.SessionId]
	sessionsLock.Unlock()
	if !exists {
		return nil, status.Errorf(codes.NotFound, "session not found")
	}

	return &pb.SessionResponse{
		SessionId:   req.SessionId,
		WebrtcOffer: peerConnection.LocalDescription().SDP,
	}, nil
}

// Chat implements a bidirectional chat stream, echoing back messages with a timestamp.
func (s *server) Chat(stream pb.SessionService_ChatServer) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Printf("Chat message from %s: %s", msg.Sender, msg.Message)
		// Add a timestamp and echo back the message.
		msg.Timestamp = time.Now().Unix()
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
}
