FROM python:3.10-alpine

WORKDIR /usr/src/app
RUN mkdir output

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt
COPY *.py ./

CMD ["python", "contrast_policy_as_code.py"]
