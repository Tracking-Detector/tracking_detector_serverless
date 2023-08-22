import os
import redis
from rq import Worker, Queue, Connection
from pymongo import MongoClient
from lib.model import get_all_available_models
import logging

logger = logging.getLogger("Model")
logging.basicConfig(level=logging.INFO)

def train_model(model_name: str, data_name: str, batch_size: int = 512, epochs: int = 10):
    print(f"Training model: {model_name} with data: {data_name}")
    for model in get_all_available_models():
        if model.model_name == model_name:
            model.train(data_name, batch_size, epochs)


def seed_models_in_db():
    logger.info("Start seeding mongodb.")
    db_client = MongoClient(os.getenv("MONGO_URI"))
    db = db_client["tracking-detector"]
    collection = db[os.getenv("MODELS_COLLECTION")]
    for model in get_all_available_models():
        collection.update_one(
            {"name": model.model_name},
            {"$set": {"name": model.model_name, "description": model.model_desc, "dims": model.input_dims}},
            upsert=True
        )
    logger.info("Finished seeding mongodb.")

if __name__ == "__main__":
    seed_models_in_db()
    with Connection(redis.StrictRedis(host='redis', port=6379, db=0)):
        worker = Worker(map(Queue, ['training']))
        worker.work()