interface Url {
  id?: number;
  url: string;
  alias: string;
  createdAt?: Date;
  updatedAt?: Date;
  expiresAt?: Date;
  owner: number;
}
