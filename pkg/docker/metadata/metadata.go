package metadata  

import (
  "context"
  // "encoding/base64"
  "fmt"  
  // "strings"
  
  // "github.com/aws/aws-sdk-go/aws/awserr"
  // "github.com/aws/aws-sdk-go/aws/session"
  // "github.com/aws/aws-sdk-go/service/ecr"
  "github.com/containers/image/transports/alltransports"
  "github.com/containers/image/types"
  "github.com/docker/cli/cli/config"
  "github.com/docker/docker/pkg/term"
  "github.com/docker/distribution/reference"

)

const (
//   awsEcrDomainSuffix = "amazonaws.com"
  kubeImageOS = "linux"
)

// TODO: provide possibility to provide credentials with parameters

// TODO: remove
func Get() { 

  // 1) check local repo || force remote
  // 2) parse remot url

  // labels := make(map[string]string)
  // _, _ = GetLabels("abcdefg:latest")
  // // fmt.Println(labels)
  // _, _ = GetLabels("abcdefg111:latest")
  // _, _ = GetLabels("abcdefg111")
  // // fmt.Println(labels)
  // labels, _ := GetLabels("docker://fedora:latest")

  labels, err := GetLabels("docker://727466838232.dkr.ecr.eu-west-1.amazonaws.com/aquila:2.0.0.2")
  // labels, _ := GetLabels("docker://starof/apache:latest")
  // labels, _ := GetLabels("docker://registry.starofservice.com/docker/apache:latest")
  if err != nil {
    fmt.Println(err.Error())
    return
  }
  fmt.Println(labels)
}


func GetLabels(name string) (map[string]string, error) {
  ref, err := alltransports.ParseImageName(name)
  if err != nil {
    panic(err)
  }

  ctx := context.Background()
  sys := &types.SystemContext{
    OSChoice: kubeImageOS,
  }

  // Trying to get metadata without authentication for public repo
  resp, err := getMetadataLabels(ctx, sys, ref)
  if err == nil {
    return resp, nil
  }

  username, password, err := getCredentials(reference.Domain(ref.DockerReference()))
  if err != nil {
    return nil, err
  }

  sys.DockerAuthConfig = &types.DockerAuthConfig{
    Username: username,
    Password: password,
  }

  resp, err = getMetadataLabels(ctx, sys, ref)
  if err != nil {
    return nil, fmt.Errorf("Unable to get image labels due to the error: %s", err.Error())
  }

  return resp, nil
}



func getMetadataLabels(ctx context.Context, sys *types.SystemContext, ref types.ImageReference) (map[string]string, error) {
  img, err := ref.NewImage(ctx, sys)
  if err != nil {
    return nil, err
  }

  imgInspect, err := img.Inspect(ctx)
  if err != nil {
    return nil, err
  }

  return imgInspect.Labels, nil
}

func getCredentials(registry string) (string, string, error) {
  // if strings.HasSuffix(registry, awsEcrDomainSuffix) {
  //   return getCredentialsAws(registry)
  // }
  _, _, stderr := term.StdStreams()
  dockerConfig := config.LoadDefaultConfigFile(stderr)
  creds, err := dockerConfig.GetAuthConfig(registry)
  if err != nil {
    return "", "", fmt.Errorf("Failed to extract docker crednetials due to the error: %s", err.Error())
  }

  if len(creds.Username) == 0 || len(creds.Password) == 0 {
    return "", "", fmt.Errorf("Got empty docker username or password")
  }

  return creds.Username, creds.Password, nil
}

// func getCredentialsAws(registry string) (string, string, error) {
//   svc := ecr.New(session.New())
  
//   input := &ecr.GetAuthorizationTokenInput{}
//   result, err := svc.GetAuthorizationToken(input)
  
//   if err != nil {
//     if aerr, ok := err.(awserr.Error); ok {
//       switch aerr.Code() {
//       case ecr.ErrCodeServerException:
//         fmt.Println(ecr.ErrCodeServerException, aerr.Error())
//       case ecr.ErrCodeInvalidParameterException:
//         fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
//       default:
//         fmt.Println(aerr.Error())
//       }
//     } else {
//       // Print the error, cast err to awserr.Error to get the Code and
//       // Message from an error.
//       fmt.Println(err.Error())
//     }
//     // return
//   }
//   authToken := result.AuthorizationData[0].AuthorizationToken
//   tokenB64Decoded, err := base64.StdEncoding.DecodeString(*authToken)
//   if err != nil {
//     return "", "", fmt.Errorf("Unable to decode string `%s` due to the error: %s", *authToken, err.Error())
//   }

//   del := strings.Index(string(tokenB64Decoded), ":")
//   if del == -1 {
//     return "", "", fmt.Errorf("Got an invalid baseauth token. The Base64 encoded string doesn't contain colon: %s", tokenB64Decoded)
//   }

//   username, password := string(tokenB64Decoded[:del]), string(tokenB64Decoded[del+1:])
//   return username, password, nil
// }
