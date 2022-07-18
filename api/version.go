//go:build !codeanalysis
// +build !codeanalysis

package api

//func (s *Server) Version(ctx context.Context, in *emptypb.Empty) (*npool.VersionResponse, error) {
//	resp, err := version.Version()
//	if err != nil {
//		logger.Sugar().Errorw("[Version] get service version error: %w", err)
//		return &npool.VersionResponse{}, status.Error(codes.Internal, "internal server error")
//	}
//	return resp, nil
//}
