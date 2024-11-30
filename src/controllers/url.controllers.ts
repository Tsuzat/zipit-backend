import {
  countUrlsByUser,
  findUrlByAlias,
  getAllUrlsByUser,
  insertUrl,
} from "../database/url.db";
import { ApiError } from "../utils/apiError";
import { asyncHandler } from "../utils/asyncHandler";
import { getXMinutesFromNow, randomSecret } from "../utils/functions";

export const createUrl = asyncHandler(async (req, res, next) => {
  let { url, alias } = req.body;
  const userId = req.user.id;
  if (!url) {
    throw new ApiError(400, "Missing required parameters");
  }
  if (!alias || alias.trim() === "") alias = randomSecret(7);
  let hasUniqueAlias = false;
  // get a unique alias
  for (let i = 0; i < 10; i++) {
    // check url exists
    const tmpUrl = await findUrlByAlias(alias);
    if (!tmpUrl || tmpUrl === null) {
      hasUniqueAlias = true;
      break;
    }
    alias = randomSecret(7);
  }
  if (!hasUniqueAlias) {
    throw new ApiError(400, "Could not create URL");
  }
  const newUrl: Url = {
    url,
    alias,
    createdAt: new Date(),
    updatedAt: new Date(),
    expiresAt: getXMinutesFromNow(5 * 24 * 60),
    owner: userId,
  };
  const newUrlId = await insertUrl(newUrl);
  newUrl.id = newUrlId;
  if (newUrlId === -1) {
    throw new ApiError(400, "Could not create URL");
  }
  res.status(200).json({ message: "URL created successfully", newUrl });
});

export const getAllUrls = asyncHandler(async (req, res, next) => {
  const userId = req.user.id;
  const urls = await getAllUrlsByUser(userId);
  res.status(200).json({ message: "URLs retrieved successfully", urls });
});

export const countUrls = asyncHandler(async (req, res, next) => {
  const userId = req.user.id;
  const count = await countUrlsByUser(userId);
  res.status(200).json({ message: "Count retrieved successfully", count });
});
