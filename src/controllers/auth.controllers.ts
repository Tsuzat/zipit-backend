import { findUser, insertUser, updateUser } from "../database/user.db";
import { UserLoginRequest, UserRegisterRequest } from "../models/auth.models";
import { User } from "../models/user.model";
import { ApiError } from "../utils/apiError";
import { ApiResponse } from "../utils/apiResponse";
import { asyncHandler } from "../utils/asyncHandler";
import { sendEmailVerification } from "../utils/emailer";
import {
  generateAccessToken,
  generateRefreshToken,
  getXMinutesFromNow,
  randomSecret,
} from "../utils/functions";
import bcrypt from "bcrypt";

export const registerUser = asyncHandler(async (req, res, next) => {
  const registerUser: UserRegisterRequest = req.body;
  const tmpUser: User | null = await findUser({ email: registerUser.email });
  if (tmpUser !== null && tmpUser.isVerified) {
    throw new ApiError(
      409,
      `User with this email already exists and is verified. Please login.`
    );
  } else if (tmpUser !== null && !tmpUser.isVerified) {
    tmpUser.name = registerUser.name;
    tmpUser.password = bcrypt.hashSync(registerUser.password, 10);
    const isUpdated = await updateUser(tmpUser);
    if (isUpdated) {
      sendEmailVerification(tmpUser);
      res
        .status(200)
        .json(
          new ApiResponse(
            200,
            "User Registered Successfully. Please check your email for verification link",
            tmpUser
          )
        );
    } else {
      throw new ApiError(500, "Error Updating User");
    }
    return;
  }
  const user: User = {
    name: registerUser.name,
    email: registerUser.email,
    password: bcrypt.hashSync(registerUser.password, 10),
    isVerified: false,
    verificationToken: randomSecret(50),
    verificationTokenExpiry: getXMinutesFromNow(30),
    tokenVersion: 1,
    isPremium: false,
    maxUrls: 0,
  };
  const userId = await insertUser(user);
  if (userId > -1) {
    sendEmailVerification(user);
    res
      .status(200)
      .json(new ApiResponse(200, "User Registered", { id: userId }));
  } else {
    console.log("Error Registering User");
    throw new ApiError(500, "Error Registering User");
  }
});

export const loginUser = asyncHandler(async (req, res, next) => {
  const loginRequest: UserLoginRequest = req.body;
  if (!loginRequest.email || !loginRequest.password) {
    throw new ApiError(400, "Missing email or password");
  }
  const user = await findUser({ email: loginRequest.email });
  if (user === null) {
    throw new ApiError(404, "User not found");
  } else if (!user.isVerified) {
    sendEmailVerification(user);
    throw new ApiError(401, "User not verified. Please verify your email");
  }
  if (!(await bcrypt.compare(loginRequest.password, user.password))) {
    throw new ApiError(401, "Invalid email or password");
  }
  console.log(user);
  const accessToken = generateAccessToken(user);
  const refreshToken = generateRefreshToken(user);
  user.refreshToken = refreshToken;
  const isUpdated = await updateUser(user);
  if (!isUpdated) {
    throw new ApiError(500, "Error Updating User");
  }
  // Send the tokens as cookies
  res.cookie("access_token", accessToken);
  res.cookie("refresh_token", refreshToken);
  res
    .status(200)
    .json(new ApiResponse(200, "Login Successful", { id: user.id }));
});

export const verifyEmail = asyncHandler(async (req, res, next) => {
  const { email, verificationToken } = req.query;
  if (!email || !verificationToken) {
    throw new ApiError(400, "Missing email or verification token");
  }
  const user = await findUser({ email });
  if (!user) {
    throw new ApiError(404, "User not found");
  }
  if (user.verificationToken !== verificationToken) {
    throw new ApiError(401, "Invalid verification token");
  }
  user.isVerified = true;
  const isUpdated = await updateUser(user);
  if (isUpdated) {
    res.status(200).json(
      new ApiResponse(200, "Email Verified Successfully", {
        isVerified: true,
      })
    );
  } else {
    throw new ApiError(500, "Error Updating User");
  }
});
