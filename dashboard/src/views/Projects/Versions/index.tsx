import SearchIcon from "@mui/icons-material/Search";
import { FormControl, Input, Sheet, Stack, Table, Typography } from "@mui/joy";
import { useMemo, useState } from "react";
import { useGetPageParams } from "../../../api/hooks/useGetPageParams";
import { useGetProject } from "../../../api/hooks/useGetProject";
import { useGetVersions } from "../../../api/hooks/useGetVersions";
import { Header } from "../../../components/Header";
import { ManageVersionModal } from "../../../components/ManageVersionModal";
import { Placeholder } from "../../../components/Placeholder";
import { TablePaginationSection } from "../../../components/TablePaginationSection";
import { Version } from "../../../interfaces/versions";
import { VersionTableRow } from "./components/VersionTableRow";

const ITEMS_PER_PAGE = 25;

export const VersionsOverview = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const { projectId } = useGetPageParams();
  const { data: project, isLoading: isProjectLoading } =
    useGetProject(projectId);
  const { data: versionsData, isLoading: isVersionsDataLoading } =
    useGetVersions(projectId);

  const filteredVersions: Version[] = useMemo(() => {
    if (
      !versionsData ||
      !versionsData.versions ||
      versionsData.versions.length === 0
    )
      return [];

    return versionsData.versions.filter((version: Version) =>
      version.name.toLowerCase().includes(searchTerm.toLowerCase()),
    );
  }, [versionsData, searchTerm]);

  const pageCount =
    filteredVersions.length > 0
      ? Math.ceil(filteredVersions.length / ITEMS_PER_PAGE)
      : 0;

  const paginatedVersions = useMemo(() => {
    if (filteredVersions.length === 0) return [];

    const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
    return filteredVersions.slice(startIndex, startIndex + ITEMS_PER_PAGE);
  }, [filteredVersions, currentPage]);

  if (isProjectLoading || isVersionsDataLoading) {
    return <Placeholder />;
  }

  return (
    <>
      <Header
        items={[
          {
            name: "Projects",
            to: "/projects",
          },
          {
            name: project?.name || "",
            to: `/projects/${projectId}/versions`,
          },
        ]}
      />

      <Stack spacing={2} direction="row">
        <FormControl sx={{ flex: 1, pb: "1rem" }}>
          <Input
            placeholder="Search for version..."
            startDecorator={<SearchIcon />}
            onChange={(e) => setSearchTerm(e.target.value)}
            value={searchTerm}
          />
        </FormControl>

        {projectId && <ManageVersionModal projectId={projectId} />}
      </Stack>

      <Sheet
        variant="outlined"
        sx={{
          display: { xs: "initial", sm: "initial" },
          width: "100%",
          borderRadius: "sm",
          flexShrink: 1,
          overflow: "auto",
          minHeight: 0,
        }}
      >
        {filteredVersions.length === 0 && (
          <Typography level="body-md" sx={{ padding: "1rem" }}>
            This project has no versions yet.
          </Typography>
        )}

        {filteredVersions.length > 0 && (
          <Table
            aria-labelledby="tableTitle"
            stickyHeader
            hoverRow
            sx={{
              "--TableCell-headBackground":
                "var(--joy-palette-background-level1)",
              "--Table-headerUnderlineThickness": "1px",
              "--TableRow-hoverBackground":
                "var(--joy-palette-background-level1)",
              "--TableCell-paddingY": "4px",
              "--TableCell-paddingX": "8px",
            }}
          >
            <thead>
              <tr>
                <th style={{ width: 240, padding: "0.7rem 1.5rem" }}>Name</th>
                <th style={{ width: 340, padding: "0.7rem 0.5rem" }}>
                  Created At
                </th>
                <th
                  style={{
                    width: 140,
                    padding: "0.7rem 4.5rem",
                    textAlign: "end",
                  }}
                ></th>
              </tr>
            </thead>

            {projectId && (
              <tbody>
                {paginatedVersions.map((version: Version) => (
                  <VersionTableRow
                    key={version.id}
                    projectId={projectId}
                    versionId={version.id}
                    versionName={version.name}
                    createdAt={version.createdAt}
                  />
                ))}
              </tbody>
            )}
          </Table>
        )}
      </Sheet>

      {filteredVersions.length > ITEMS_PER_PAGE && (
        <TablePaginationSection
          currentPage={currentPage}
          pageCount={pageCount}
          onPageChange={setCurrentPage}
        />
      )}
    </>
  );
};
