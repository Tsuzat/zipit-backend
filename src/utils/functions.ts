import crypto from "crypto";
import jwt from "jsonwebtoken";
import { User } from "../models/user.model";
import {
  ACCESS_TOKEN_EXPIRY,
  ACCESS_TOKEN_SECRET,
  REFRESH_TOKEN_SECRET,
  REFRESH_TOKEN_EXPIRY,
} from "../constants";

const randomSecret = (length: number = 36): string =>
  crypto.randomBytes(length).toString("hex");

const getXMinutesFromNow = (x: number): Date => {
  const now = new Date();
  now.setTime(now.getTime() + x * 60 * 1000);
  return now;
};

const generateAccessToken = (user: User): string => {
  return jwt.sign(
    {
      id: user.id,
      email: user.email,
      name: user.name,
      token_version: user.tokenVersion,
    },
    ACCESS_TOKEN_SECRET,
    {
      expiresIn: ACCESS_TOKEN_EXPIRY,
    }
  );
};

const generateRefreshToken = (user: User): string => {
  return jwt.sign(
    {
      id: user.id,
    },
    REFRESH_TOKEN_SECRET,
    {
      expiresIn: REFRESH_TOKEN_EXPIRY,
    }
  );
};

export {
  randomSecret,
  getXMinutesFromNow,
  generateAccessToken,
  generateRefreshToken,
};
