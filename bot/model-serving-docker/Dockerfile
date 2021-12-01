# syntax=docker/dockerfile:1
FROM tensorflow/tensorflow:latest
WORKDIR /home
RUN echo "deb [arch=amd64] http://storage.googleapis.com/tensorflow-serving-apt stable tensorflow-model-server tensorflow-model-server-universal" | \
    tee /etc/apt/sources.list.d/tensorflow-serving.list && \
    curl https://storage.googleapis.com/tensorflow-serving-apt/tensorflow-serving.release.pub.gpg | apt-key add - && \
    apt-get update && apt-get install tensorflow-model-server
EXPOSE 7070
CMD "/home/tf/serve.sh"