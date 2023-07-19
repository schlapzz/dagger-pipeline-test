FROM golang:1.20


WORKDIR /src

RUN wget -O oc.tar https://downloads-openshift-console.ocp-internal.cloudscale.puzzle.ch/amd64/linux/oc.tar && tar -xvf oc.tar && rm oc.tar && chmod +x ./oc && mv ./oc  /usr/local/bin/oc
RUN wget -O dagger.tar.gz https://github.com/dagger/dagger/releases/download/v0.6.3/dagger_v0.6.3_linux_amd64.tar.gz && tar -xvf dagger.tar.gz && rm dagger.tar.gz && chmod +x ./dagger && mv ./dagger  /usr/local/bin/dagger

RUN dagger --help