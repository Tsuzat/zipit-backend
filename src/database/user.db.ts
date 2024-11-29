import { error } from "console";
import { User } from "../models/user.model";
import db from "./db";
import { users } from "./schemas";
import { eq, and } from "drizzle-orm";

async function insertUser(user: User): Promise<number> {
  const result = await db
    .insert(users)
    .values(user)
    .returning({ id: users.id });
  if (result.length === 0) {
    return -1;
  }
  return result[0].id;
}

type FindUserParams = {
  id?: number;
  email?: string;
};

/**
 * Function to find a user by either ID or email
 * @param params Find user parameters - either 'id' or 'email' must be provided
 * @returns Returns the user object if found, otherwise returns null
 */
async function findUser(params: FindUserParams): Promise<User | null> {
  const { id, email } = params;

  // Ensure at least one parameter is provided
  if (!id && !email) {
    throw new Error("Either 'id' or 'email' must be provided.");
  }

  // Build the where clause dynamically
  const conditions = [];
  if (id) conditions.push(eq(users.id, id));
  if (email) conditions.push(eq(users.email, email));

  const query = db
    .select()
    .from(users)
    .where(and(...conditions)) // Combine conditions with AND
    .limit(1); // Limit the query to one user

  const [user] = await query; // Destructure to get the single user
  return user || null; // Return null if no user found
}

async function updateUser(user: User): Promise<boolean> {
  const result = await db
    .update(users)
    .set(user)
    .where(eq(users.id, user.id))
    .returning({ id: users.id });
  if (result.length === 0) {
    return false;
  }
  return true;
}

export { insertUser, findUser, updateUser };
