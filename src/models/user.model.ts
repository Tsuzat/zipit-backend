export interface User {
  id?: number;
  name: string;
  email: string;
  password: string;
  profileImage?: string;
  isVerified: boolean;
  refreshToken?: string;
  verificationToken?: string;
  verificationTokenExpiry?: Date;
  tokenVersion: number;
  createdAt?: Date;
  updatedAt?: Date;
  isPremium: boolean;
  maxUrls: number;
}
