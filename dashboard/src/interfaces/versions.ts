export type GetVersionsResponse = {
  versions: Version[];
};

export type Version = {
  id: number;
  name: string;
  createdAt: string; // ISO 8601
};
