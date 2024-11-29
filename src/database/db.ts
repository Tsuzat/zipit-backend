import { DB_URL } from "../constants";
import { drizzle } from "drizzle-orm/node-postgres";

const db = drizzle(DB_URL, {});

export default db;
