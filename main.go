package main

import (
	"context"
	"fmt"
	"log"
	"io"
	"io/ioutil"
	"os"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func main() {
	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Println("Failed to create client!")
	}

	file := "bthorinterview.wav"

	err  = sendAudio(os.Stdout, client, file)
	if err != nil {
		log.Println(err.Error())
	}


}

func sendAudio(w io.Writer, client *speech.Client, file string) error {
	ctx := context.Background()
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	request := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:	speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode: "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}

	op, err := client.LongRunningRecognize(ctx, request)
	if err != nil {
		return err
	}


	response, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	for _, results := range response.Results {
		for _, alt := range results.Alternatives {
			fmt.Fprintf(w, "\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}

	return nil
}
