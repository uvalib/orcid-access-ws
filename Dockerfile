FROM alpine:3.7

# update the packages
RUN apk update && apk upgrade && apk add bash tzdata ca-certificates && rm -fr /var/cache/apk/*

# Create the run user and group
RUN addgroup webservice && adduser webservice -G webservice -D

# set the timezone appropriatly
ENV TZ=UTC
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Specify home 
ENV APP_HOME /orcid-access-ws
WORKDIR $APP_HOME

# Create necessary directories
RUN mkdir -p $APP_HOME/scripts $APP_HOME/bin $APP_HOME/data
RUN chown -R webservice $APP_HOME && chgrp -R webservice $APP_HOME

# Specify the user
USER webservice

# port and run command
EXPOSE 8080
CMD scripts/entry.sh

# Move in necessary assets
COPY data/container_bash_profile /home/webservice/.profile
COPY scripts/entry.sh $APP_HOME/scripts/entry.sh
COPY data/work-activity-template.xml $APP_HOME/data/work-activity-template.xml
COPY bin/orcid-access-ws.linux $APP_HOME/bin/orcid-access-ws

# Add the build tag
COPY buildtag.* $APP_HOME/
