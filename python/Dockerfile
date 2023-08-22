FROM toddstavish/tensorflow.cpu
RUN pip install tensorflowjs
# Set the working directory in the container
WORKDIR /app

# Copy the requirements file
COPY ./python /app
RUN pip install --no-cache-dir --upgrade protobuf

# Install the Python dependencies
RUN pip install --no-cache-dir -r requirements.txt


# Set the command to run the FastAPI application using Uvicorn
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8081"]