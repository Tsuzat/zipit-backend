import {
  pgTable,
  serial,
  text,
  boolean,
  integer,
  timestamp,
} from "drizzle-orm/pg-core";

export const users = pgTable("users", {
  id: serial("id").primaryKey(),
  name: text("name").notNull(),
  email: text("email").notNull(),
  password: text("password").notNull(),
  profileImage: text("profile_image"),
  isVerified: boolean("is_verified").default(false),
  refreshToken: text("refresh_token"),
  verificationToken: text("verification_token"),
  verificationTokenExpiry: timestamp("verification_token_expiry"),
  tokenVersion: integer("token_version").default(0),
  createdAt: timestamp("created_at").defaultNow(),
  updatedAt: timestamp("updated_at").defaultNow(),
  isPremium: boolean("is_premium").default(false),
  maxUrls: integer("max_urls").default(0),
});

export const urls = pgTable("urls", {
  id: serial("id").primaryKey(),
  url: text("url").notNull(),
  alias: text("alias").notNull(),
  createdAt: timestamp("created_at").defaultNow(),
  updatedAt: timestamp("updated_at").defaultNow(),
  expiresAt: timestamp("expires_at"),
  owner: integer("owner").notNull(),
});
