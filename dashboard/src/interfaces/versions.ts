export interface GetVersionsResponse {
  versions: Version[];
}

export interface Version {
  id: number;
  name: string;
  createdAt: string; // ISO 8601
}
