export type GetProjectsResponse = {
  projects: Project[];
}

export type postProjectRequest = {
  name: string;
}

export type Project = {
  id: number;
  name: string;
  numberOfVersions: number;
  createdAt: string; // ISO 8601
}
