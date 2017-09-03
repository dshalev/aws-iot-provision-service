package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iot/iotiface"
	"fmt"
	"os"
	"encoding/json"
	"github.com/dshalev2/aws-iot-provision-service/model"
)

var (
	svc iotiface.IoTAPI
	thingType string
)

func HandleProvision(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	thingName := vars["thingName"]

	svc = iot.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

	ctResp, err := svc.CreateThing(
		&iot.CreateThingInput{
			ThingName:        aws.String(thingName),
			AttributePayload: buildAttributes(),
		})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create thing: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Failed to create thing"))
	}

	fmt.Println("Thing ARN: ", *ctResp.ThingArn)

	resp, err := svc.CreateKeysAndCertificate(&iot.CreateKeysAndCertificateInput{
		SetAsActive: aws.Bool(true),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create certificate: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Failed to create certificate"))
	}

	tconfig := model.NewThingConfig(resp)

	jData, err := json.Marshal(tconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal json: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Failed to marshal json"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}


func buildAttributes() *iot.AttributePayload {

	ap := &iot.AttributePayload{
		Attributes: map[string]*string{},
	}

	if thingType != "" {
		ap.Attributes["Type"] = aws.String(thingType)
	}

	return ap
}