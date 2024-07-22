import SearchIcon from "@mui/icons-material/Search";
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  Sheet,
  Table,
  Typography,
} from "@mui/joy";
import { useEffect, useState } from "react";
import { useParams } from "react-router";
import { useGetProject } from "../../../api/hooks/useGetProject";
import { useGetVersions } from "../../../api/hooks/useGetVersions";
import { PaginationSection } from "../../../components/TablePaginationSection";
import { GetVersionsResponse } from "../../../interfaces/versions";
import { ProjectTableRow } from "../components";

export const VersionsOverview = () => {
  const [searchBar, setSearchBar] = useState("");
  const { projectId } = useParams();
  const { data: project } = useGetProject(projectId);
  const { data: versionsData } = useGetVersions(projectId);
  const [versionsList, setVersionsList] = useState<GetVersionsResponse>();

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

      <Button sx={{ mb: "0.5rem", backgroundColor: '#0078ff' }}>Add New Version</Button>

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

          {versionsList && (
            <tbody>
              {versionsList.versions.map((version) => (
                <ProjectTableRow
                  key={version.id}
                  id={version.id}
                  name={version.name}
                  createdAt={version.createdAt}
                />
              ))}
            </tbody>
          )}
        </Table>
      </Sheet>

      <PaginationSection />
    </>
  );
};
