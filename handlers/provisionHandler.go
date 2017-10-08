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
	decoder := json.NewDecoder(r.Body)
	var b model.CsrConfig
	_ = decoder.Decode(&b)


	resp, err := svc.CreateCertificateFromCsr(&iot.CreateCertificateFromCsrInput{
		CertificateSigningRequest: aws.String(b.CsrText),
		SetAsActive: aws.Bool(true),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create certificate: %v %v\n", b, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Failed to create certificate"))
	}

	presp, err := svc.GetPolicy(&iot.GetPolicyInput{
		PolicyName:     aws.String("PubSubToAnyTopic"),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get policy: %v\n", err)
		w.Write([]byte("500 - Failed to get policy"))
	}


	_, err = svc.AttachPrincipalPolicy(&iot.AttachPrincipalPolicyInput{
		PolicyName: presp.PolicyName,
		Principal:  resp.CertificateArn,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to attach policy: %v\n", err)
		w.Write([]byte("500 - Failed to get policy"))
	}

	_, err = svc.AttachThingPrincipal(&iot.AttachThingPrincipalInput{
		Principal: resp.CertificateArn,
		ThingName: aws.String(thingName),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to attach thing to cert: %v\n", err)
		w.Write([]byte("500 - Failed to get policy"))
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
