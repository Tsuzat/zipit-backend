import { Router } from "express";
import {
  loginUser,
  registerUser,
  verifyEmail,
} from "../controllers/auth.controllers";

const authRouter = Router();

authRouter.route("/signup").post(registerUser);
authRouter.route("/verify-email").get(verifyEmail);
authRouter.route("/login").post(loginUser);

export { authRouter };
