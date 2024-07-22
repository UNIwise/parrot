export interface GetProjectsResponse {
  projects: Project[];
}

export interface Project {
  id: number;
  name: string;
  numberOfVersions: number;
  createdAt: string; // ISO 8601
}
