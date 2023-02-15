package git

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	github "github.com/google/go-github/v50/github"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const GH_Token = "ghp_"
const repo = "blueprints"
const owner = "ntnguyen-dcn"
const GithubBranch = "main"

type GitClient struct {
	repo   string
	owner  string
	client *github.Client
	logger logr.Logger
}

func NewClient(repo, owner, token string, ctx context.Context) (GitClient, error) {
	var client GitClient
	client.repo = repo
	client.owner = owner
	if token == "" {
		token = GH_Token
	}
	client.logger = ctrl.Log.WithName("Github Module: ")
	client.client = github.NewTokenClient(ctx, token)
	return client, nil
}
func (client *GitClient) CommitNewFile(fileName, branch, folder string, content []byte) (*github.RepositoryContentResponse, error) {
	commitMessage := "Commit a new file: " + fileName
	path := folder + fileName
	commitOptions := &github.RepositoryContentFileOptions{
		Message: &commitMessage,
		Content: []byte(content),
		Branch:  &branch,
	}
	contentRsp, _, err1 := client.client.Repositories.CreateFile(context.Background(), client.owner, client.repo, path, commitOptions)
	if err1 != nil {
		client.logger.Error(err1, "Commit file failed")
	}
	return contentRsp, err1
}

func (client *GitClient) GetFileHistory(fileName string, folder string) ([]*github.RepositoryCommit, error) {
	filePath := folder + fileName

	commits, _, err := client.client.Repositories.ListCommits(context.Background(), owner, repo, &github.CommitsListOptions{Path: filePath})
	if err != nil {
		client.logger.Error(err, "Get History file failed")
	}
	return commits, err
}

func (client *GitClient) GetContentOfFile(ctx context.Context, fileName, folder string, ref string) (*github.RepositoryContent, error) {
	path := folder + fileName
	var options github.RepositoryContentGetOptions
	if ref != "" {
		options.Ref = ref
	}
	content, _, _, err := client.client.Repositories.GetContents(ctx, client.owner, client.repo, path, &options)
	if err != nil {
		client.logger.Error(err, "Get Content of file failed")
	}
	return content, err
}

func (client *GitClient) UpdateFile(fileName, folder string, content []byte) error {
	filePath := folder + fileName
	// Retrieve the current content of the file
	fileContent, _, _, err := client.client.Repositories.GetContents(context.TODO(), owner, repo, filePath, &github.RepositoryContentGetOptions{})
	if err != nil {
		client.logger.Error(err, "Get File Info failed")
		return err
	}
	// Update the file in the repository
	updateOpts := &github.RepositoryContentFileOptions{
		Message: github.String("Update file content"),
		Content: content,
		SHA:     fileContent.SHA,
	}
	_, _, err1 := client.client.Repositories.UpdateFile(context.TODO(), owner, repo, filePath, updateOpts)
	if err1 != nil {
		client.logger.Error(err1, "Update File content failed")
	}
	return err1
}

// func (client *GitClient) UpdateRevision(revision string, fileName string, folder string) error {
// 	path := folder + fileName + ".revision"

// }

// =====================================================
// Test client API part

func CommitFile(filename string, content []byte) error {
	logger := log.Log.WithName("Git modules")
	logger.Info("Using git module...\n")
	ctx := context.Background()
	client := github.NewTokenClient(ctx, GH_Token)
	user, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		logger.Error(err, "Error when auth with Github server\n")
		return err
	} else {
		logger.Info("Authenticated with Github", "user", user.Name, "Response Status", resp.StatusCode)
	}
	// fileContentEncoded := base64.StdEncoding.EncodeToString(content)
	// Create a new commit with the file
	logger.Info("Committing file...\n")
	commitMessage := "Commit " + filename
	path := filename + ".yaml"
	branch := "main"
	commitOptions := &github.RepositoryContentFileOptions{
		Message: &commitMessage,
		Content: []byte(content),
		Branch:  &branch,
	}
	contentRsp, rsp, err1 := client.Repositories.CreateFile(context.Background(), owner, repo, path, commitOptions)
	if err1 != nil {
		logger.Error(err1, "Error when committing file\n")
		return err1
	} else {
		logger.Info("Commited", *(contentRsp.SHA), rsp.StatusCode)
	}
	return nil
}

func GetFileHistory(filePath string) error {
	filePath = filePath + ".yaml"
	logger := log.Log.WithName("Git modules")
	logger.Info("Using git module...Get History\n")
	ctx := context.Background()
	client := github.NewTokenClient(ctx, GH_Token)
	// Retrieve the commit history for the file
	commits, _, err := client.Repositories.ListCommits(context.Background(), owner, repo, &github.CommitsListOptions{Path: filePath})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving commit history: %s\n", err.Error())
		return err
	}
	logger.Info("List history:\n", "History", commits)
	return nil

}
