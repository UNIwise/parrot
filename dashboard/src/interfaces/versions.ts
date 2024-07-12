export interface GetVersionsResponse {
  versions: Version[];
}

export interface Version {
  id: number;
  name: string;
  createdAt: string; // ISO 8601
}
    // background-color: var(--variant-solidBg, var(--joy-palette-primary-solidBg, var(--joy-palette-primary-500, #0B6BCB)));
