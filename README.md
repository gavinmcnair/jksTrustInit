# jksTrustInit

![GitHub](https://img.shields.io/github/license/gavinmcnair/jkstrustinit)
![CircleCI](https://img.shields.io/circleci/build/github/gavinmcnair/jksTrustInit/main?token=aab7daba901f49034a2fb9f61895b61114b13de9)


## Problem statement

Have you got a Java application which needs to deliver a `JKS` file but you only have a standard pem encoded Key and Certificate?

jksTrustInit is an `initContainer` which takes creds from either local files or environment variables and writes out a Java Keystore (JKS) file to an emptyDir which can be shared with the main container

| Environment Variable  | Default  | Description  |
|---|---|---|
| PASSWORD  | password  | The password used for the keystore|
| FILE_MODE  | false | If to use the env vars or files  |
| KEY  |  NA | Public Key environment variable |
| CERTIFICATE  |  NA | Certificate environment variable  |
| KEY_FILE  |  NA |  Public Key file |
| CERTIFICATE_FILE  | NA  | Certificate file  |
| OUTPUT_FILE  | /var/run/secrets/truststore.jks  | The filename used to write the file out |

## How to use in Kubernetes

We can supply the PEM encoded `key` and `certificate` either within the environment variable or as files mounted upon the filesystem. Both of which can be sourced with secrets or configmaps as appropriate. When using files you need to set `FILE_MODE` to `true`

The init container will start and write the output file to the `OUTPUT_FILE` path.

This is then available to the target JVM.

### Example pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: KafkaClient
spec:
  initContainers:
    - name: jksTrustInit
      image: gavinmcnair/jksTrustInit:1.0.0
      env:
        - name: KEY
          value: "pem encoded key"
        - name: CERTIFICATE
          value: "pem encoded cert"
      volumeMounts:
        - mountPath: /var/run/secrets
          name: kafkasecrets
  containers:
    - name: kafkaclient
      image: kafkaclient:1.0.0
      env:
        - name: JAVA_JKS_FILE
          value: "/var/run/secrets/truststore.jks"
        - name: JAVA_JKS_PASSWORD
          value: "password"
      volumeMounts:
        - mountPath: /var/run/secrets
          name: kafkasecrets
  volumes:
    - emptyDir: {}
      name: kafkasecrets

```

## Motivation

To do this in the conventional way you need to use an insecure Java container which is large and execute multiple java keystore commands. This container is a single binary on a scratch container.

It should be both quick and reliable.
