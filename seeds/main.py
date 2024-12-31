import random
from faker import Faker
from datetime import datetime
from dataclasses import dataclass
import os
from dotenv import load_dotenv
from pymongo import MongoClient
import logging

logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s"
)

dotenv_path = os.path.join(os.path.dirname(os.path.dirname(__file__)), ".env")

load_dotenv(dotenv_path)

fake = Faker(["vi_VN", "en_US", "ja_JP"])


@dataclass
class User:
    username: str
    hash_password: str
    avatar: str
    email: str
    first_name: str
    last_name: str
    role: list[str]
    created_at: datetime
    updated_at: datetime

    def to_dict(self):
        return {
            "username": self.username,
            "hash_password": self.hash_password,
            "avatar": self.avatar,
            "email": self.email,
            "first_name": self.first_name,
            "last_name": self.last_name,
            "role": self.role,
            "created_at": self.created_at,
            "updated_at": self.updated_at,
        }


def generate_user() -> User:
    roles = ["user", "manager", "admin"]
    role_probabilities = [0.7, 0.2, 0.1]

    num_roles = random.randint(1, len(roles))
    assigned_roles = random.choices(roles, weights=role_probabilities, k=num_roles)

    return User(
        username=fake.user_name(),
        hash_password=fake.sha256(),
        avatar=fake.image_url(),
        email=fake.email(),
        first_name=fake.first_name(),
        last_name=fake.last_name(),
        role=list(set(assigned_roles)),
        created_at=fake.date_time_this_year(),
        updated_at=fake.date_time_this_year(),
    )


def connection_db():
    db_connection = os.getenv("DB_CONNECTION")
    try:
        client = MongoClient(db_connection)
        logging.info("Connected to MongoDB successfully! ✅")
        db = client["db"]
        return db
    except Exception as e:
        logging.error(f"Error connecting to MongoDB: {e}")


def seed_users(db, num: int) -> None:
    users_to_insert = []

    for i in range(num):
        user = generate_user()
        users_to_insert.append(user.to_dict())

    db["users"].insert_many(users_to_insert)

    logging.info("Seed users inserted successfully! 🚀")


def main() -> None:
    db = connection_db()

    # seed_users(db,10000)

    logging.info(db["users"].count_documents({}))


if __name__ == "__main__":
    main()
