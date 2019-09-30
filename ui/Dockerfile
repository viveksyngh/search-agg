FROM python:2.7-alpine

WORKDIR /home/app/

RUN mkdir -p /home/app/python 
ENV PATH=$PATH:/home/app/.local/bin:/home/app/python/bin/
ENV PYTHONPATH=$PYTHONPATH:/home/app/python

EXPOSE 8000

COPY . .

CMD ["python", "-m", "SimpleHTTPServer", "8000"]