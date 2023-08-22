FROM tensorflow/tensorflow
# Set the working directory in the container
WORKDIR /app

# Copy the requirements file
COPY ./python /app
RUN pip install --no-cache-dir --upgrade protobuf

# Install the Python dependencies
RUN pip install --no-cache-dir -r requirements.txt


# Set the command to run the FastAPI application using Uvicorn
CMD [ "python", "worker.py" ]