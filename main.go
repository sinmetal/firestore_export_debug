package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/morikuni/failure"
	"google.golang.org/api/firestore/v1"
)

func main() {
	fpID, outputPath := getArgs()
	fmt.Printf("FIRESTORE_PROJECT:%s, OUTPUT_PATH:%s\n", fpID, outputPath)

	ctx := context.Background()
	ope, err := export(ctx, fpID, outputPath)
	if err != nil {
		fmt.Printf("failed firestore export : %v\n", err)
		os.Exit(1)
	}
	if ope.HTTPStatusCode != http.StatusOK {
		j, err := json.Marshal(ope)
		if err != nil {
			fmt.Printf("failed json.Marshal : %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("failed firestore export : api response is %s", string(j))
		os.Exit(1)
	}

	fmt.Printf("done. JobID is %s\n", ope.Name)
	os.Exit(0)
}

func getArgs() (string, string) {
	var (
		flagFPID       = flag.String("project", "", "firestore export project id")
		flagOutputPath = flag.String("output", "", "output gcs path")
	)
	flag.Parse()

	envFPID := os.Getenv("FIRESTORE_PROJECT")
	envOutputPath := os.Getenv("OUTPUT_PATH")

	if len(*flagFPID) < 1 {
		*flagFPID = envFPID
	}
	if len(*flagOutputPath) < 1 {
		*flagOutputPath = envOutputPath
	}
	return *flagFPID, *flagOutputPath
}

func export(ctx context.Context, projectID string, outputPath string) (*firestore.GoogleLongrunningOperation, error) {
	f, err := firestore.NewService(ctx)
	if err != nil {
		return nil, failure.Wrap(err, failure.Messagef("failed firestore.NewService"))
	}
	name := fmt.Sprintf("projects/%s/databases/(default)", projectID)
	return f.Projects.Databases.ExportDocuments(name, &firestore.GoogleFirestoreAdminV1ExportDocumentsRequest{
		OutputUriPrefix: outputPath,
	}).Do()
}
