import { Router } from "express";
import {
  loginUser,
  Me,
  registerUser,
  verifyEmail,
} from "../controllers/auth.controllers";
import { authenticate } from "../middleware/auth.middleware";

const authRouter = Router();

authRouter.route("/signup").post(registerUser);
authRouter.route("/verify-email").get(verifyEmail);
authRouter.route("/login").post(loginUser);
authRouter.route("/me").get(authenticate, Me);

export { authRouter };
