import { Router } from "express";
import {
  countUrls,
  createUrl,
  getAllUrls,
} from "../controllers/url.controllers";
import { authenticate } from "../middleware/auth.middleware";

const urlRouter = Router();

urlRouter.route("/").post(authenticate, createUrl);
urlRouter.route("/").get(authenticate, getAllUrls);
urlRouter.route("/count").get(authenticate, countUrls);

export { urlRouter };
