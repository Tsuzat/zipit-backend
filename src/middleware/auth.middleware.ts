import { ACCESS_TOKEN_SECRET } from "../constants";
import { findUser } from "../database/user.db";
import { ApiError } from "../utils/apiError";
import { asyncHandler } from "../utils/asyncHandler";
import jwt, { TokenExpiredError } from "jsonwebtoken";

export const authenticate = asyncHandler(async (req, res, next) => {
  try {
    const token =
      req.cookies?.access_token ||
      req.headers?.authorization?.replace("Bearer ", "");
    if (!token) throw new ApiError(401, "Unauthorized Access");
    const decodeToken = jwt.verify(token, ACCESS_TOKEN_SECRET);
    const user = await findUser({
      // @ts-ignore
      id: decodeToken.id,
      // @ts-ignore
      email: decodeToken.email,
    });

    if (user === null) throw new ApiError(401, "Unauthorized User");
    req.user = user;
    next();
  } catch (error) {
    if (error instanceof TokenExpiredError)
      throw new ApiError(403, "Token Expired");
    console.log(error);
    throw new ApiError(401, "Unauthorized Access");
  }
});
