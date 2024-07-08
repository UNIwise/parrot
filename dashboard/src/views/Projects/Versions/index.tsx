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
import { mockedVersionsResponse } from "../../../api/mocks/versions.mock";
import { PaginationSection } from "../../../components/TablePaginationSection";
import { TableRow } from "../../../components/TableRow";
import { GetVersionsResponse } from "../../../interfaces/versions";

export const VersionsOverview = () => {
  const [searchBar, setSearchBar] = useState("");
  //TODO: replace mocked data with the response from the API when react query hooks are implemented
  const versions = mockedVersionsResponse.versions;
  const [versionsList, setVersionsList] = useState<GetVersionsResponse>(
    mockedVersionsResponse,
  );

  useEffect(() => {
    if (searchBar === "") {
      return;
    }
    setVersionsList(versionsList);
  }, [versionsList, searchBar]);

  const versionSearchHandle = (search: string) => {
    const filteredVersions = versions.filter((version) =>
      version.name.toLowerCase().includes(search.toLowerCase()),
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
        }}
      >
        <Typography
          level="h2"
          component="h1"
          style={{
            alignSelf: "center",
            fontSize: "3rem",
            marginRight: "1.5rem",
          }}
        >
          Versions
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

      <Button sx={{ mb: "0.5rem" }}>Add New Version</Button>

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

          <tbody>
            {versionsList.versions.map((version) => (
              <TableRow
                key={version.id}
                name={version.name}
                createdAt={version.createdAt}
              />
            ))}
          </tbody>
        </Table>
      </Sheet>

      <PaginationSection />
    </>
  );
};
