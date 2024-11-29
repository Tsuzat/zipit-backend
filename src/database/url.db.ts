import { eq } from "drizzle-orm";
import "../models/url.model";
import db from "./db";
import { urls } from "./schemas";

export const insertUrl = async (url: Url): Promise<number> => {
  const results = await db.insert(urls).values(url).returning({ id: urls.id });
  if (results.length === 0) {
    return -1;
  }
  return results[0].id;
};

export const findUrlByAlias = async (alias: string): Promise<Url | null> => {
  const query = db.select().from(urls).where(eq(urls.alias, alias)).limit(1);
  const [url] = await query;
  return url || null;
};

export const updateUrl = async (url: Url): Promise<boolean> => {
  const result = await db
    .update(urls)
    .set(url)
    .where(eq(urls.id, url.id))
    .returning({ id: urls.id });
  if (result.length === 0) {
    return false;
  }
  return true;
};

export const deleteUrl = async (id: number): Promise<boolean> => {
  const result = await db
    .delete(urls)
    .where(eq(urls.id, id))
    .returning({ id: urls.id });
  if (result.length === 0) {
    return false;
  }
  return true;
};

export const getAllUrlsByUser = async (userId: number): Promise<Url[]> => {
  const query = db.select().from(urls).where(eq(urls.owner, userId));
  return await query;
};
