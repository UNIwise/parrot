import SearchIcon from "@mui/icons-material/Search";
import {
  Box,
  FormControl,
  FormLabel,
  Input,
  Sheet,
  Table,
  Typography
} from "@mui/joy";
import { useEffect, useState } from "react";
import { useGetPageParams } from "../../../api/hooks/useGetPageParams";
import { useGetProject } from "../../../api/hooks/useGetProject";
import { useGetVersions } from "../../../api/hooks/useGetVersions";
import { TablePaginationSection } from "../../../components/TablePaginationSection";

import { ManageVersionModal } from "../../../components/Modal";
import { GetVersionsResponse, Version } from "../../../interfaces/versions";
import { VersionTableRow } from "./components";

const ITEMS_PER_PAGE = 20;

export const VersionsOverview = () => {
  const [searchBar, setSearchBar] = useState("");
  const { projectId } = useGetPageParams();
  const { data: project } = useGetProject(projectId);
  const { data: versionsData } = useGetVersions(projectId);
  const [versionsList, setVersionsList] = useState<GetVersionsResponse>();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageCount, setPageCount] = useState(1);
  const [paginatedVersions, setPaginatedVersions] = useState<Version[]>();

  useEffect(() => {
    if (!versionsData) return;

    setVersionsList(versionsData);
  }, [versionsData]);

  const versionSearchHandle = (versionName: string) => {
    if (!versionsData) return;

    const filteredVersions = versionsData.versions.filter((version) =>
      version.name.toLowerCase().includes(versionName.toLowerCase()),
    );
    setVersionsList({ versions: filteredVersions });
  };

  useEffect(() => {
    if (!versionsList || !versionsList.versions
    ) return;

    const pageCount = Math.ceil(versionsList.versions.length / ITEMS_PER_PAGE);
    const paginatedVersions = versionsList.versions.slice(
      (currentPage - 1) * ITEMS_PER_PAGE,
      currentPage * ITEMS_PER_PAGE
    );
    setPageCount(pageCount);
    setPaginatedVersions(paginatedVersions);
  }, [versionsList, currentPage]);

  const handlePageChange = (newPage: number) => {
    setCurrentPage(newPage);
  };

  return (
    <>
      <Box
        sx={{
          borderRadius: "sm",
          py: 0.5,
          display: { xs: "none", sm: "flex" },
          flexWrap: "wrap",
          gap: 1.5,
          "& > *": {
            minWidth: { xs: "120px", md: "160px" },
          },
          flexDirection: 'column',
        }}
      >
        <Typography
          level="h2"
          component="h1"
          sx={{
            alignSelf: "center",
            fontSize: "3rem",
            marginRight: "1.5rem",
            backgroundColor: '#0078ff',
            p: '1rem 5rem',
            borderRadius: "sm",
          }}
        >
          {project?.name} versions
        </Typography>

        <FormControl sx={{ flex: 1, pb: "1.1rem" }} size="sm">
          <FormLabel>Search for version</FormLabel>

          <Input
            size="sm"
            placeholder="Look up for version name"
            startDecorator={<SearchIcon />}
            onChange={(e) => {
              setSearchBar(e.target.value);
              versionSearchHandle(e.target.value);
            }}
            value={searchBar}
          />
        </FormControl>
      </Box>

      {projectId &&
        <ManageVersionModal projectId={projectId} />
      }

      <Sheet
        variant="outlined"
        sx={{
          display: { xs: "none", sm: "initial" },
          width: "100%",
          borderRadius: "sm",
          flexShrink: 1,
          overflow: "auto",
          minHeight: 0,
        }}
      >
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
              >
                Delete
              </th>
            </tr>
          </thead>

          {paginatedVersions && projectId && (
            <tbody>
              {paginatedVersions.map((version) => (
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
      </Sheet>

      <TablePaginationSection
        currentPage={currentPage}
        pageCount={pageCount}
        onPageChange={handlePageChange}
      />
    </>
  );
};
