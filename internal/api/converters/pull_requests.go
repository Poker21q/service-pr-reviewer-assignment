package converters

import (
	"service-pr-reviewer-assignment/internal/generated/api/dto"
	"service-pr-reviewer-assignment/internal/service/entities"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
)

func UserReviewRequestsToDTO(userID uuid.UUID, prs []entities.PullRequest) dto.UserReviewResponse {
	prDTOs := make([]dto.PullRequestShort, 0, len(prs))
	for _, pr := range prs {
		prDTOs = append(prDTOs, PullRequestToShortDTO(pr))
	}

	return dto.UserReviewResponse{
		UserId:       userID,
		PullRequests: prDTOs,
	}
}

func PullRequestToShortDTO(pr entities.PullRequest) dto.PullRequestShort {
	return dto.PullRequestShort{
		AuthorId:        pr.AuthorID,
		PullRequestId:   pr.ID,
		PullRequestName: pr.Name,
		Status:          PullRequestStatusToShortDTO(pr.Status),
	}
}

func PullRequestStatusToShortDTO(status entities.PullRequestStatus) dto.PullRequestShortStatus {
	switch status {
	case entities.PullRequestStatusOpen:
		return dto.PullRequestShortStatusOPEN
	case entities.PullRequestStatusMerged:
		return dto.PullRequestShortStatusMERGED
	default:
		return "UNKNOWN"
	}
}

func PullRequestToDTO(pr *entities.PullRequest, reviewerIDs []uuid.UUID) dto.PullRequest {
	return dto.PullRequest{
		AssignedReviewers: reviewerIDs,
		AuthorId:          pr.AuthorID,
		CreatedAt:         pointer.To(pr.CreatedAt),
		MergedAt:          pr.MergedAt,
		PullRequestId:     pr.ID,
		PullRequestName:   pr.Name,
		Status:            PullRequestStatusToDTO(pr.Status),
	}
}

func PullRequestStatusToDTO(status entities.PullRequestStatus) dto.PullRequestStatus {
	switch status {
	case entities.PullRequestStatusOpen:
		return dto.PullRequestStatusOPEN
	case entities.PullRequestStatusMerged:
		return dto.PullRequestStatusMERGED
	default:
		return "UNKNOWN"
	}
}

func ReassignResponseToDTO(pr *entities.PullRequest, reviewerIDs []uuid.UUID, newReviewerID uuid.UUID) dto.ReassignPullRequestResponse {
	return dto.ReassignPullRequestResponse{
		Pr:         PullRequestToDTO(pr, reviewerIDs),
		ReplacedBy: newReviewerID,
	}
}
