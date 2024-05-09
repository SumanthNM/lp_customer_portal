FROM ubuntu:latest
RUN apt-get update
RUN apt-get install ca-certificates -y
ENV CHASSIS_HOME=.
COPY ./server .
CMD [ "./main" ]