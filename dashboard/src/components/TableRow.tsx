import { Delete } from "@mui/icons-material";
import { Button, Typography } from "@mui/joy";
import { FC } from "react";
import { Link, useParams } from "react-router-dom";
import { useDeleteVersion } from "../api/hooks/useDeleteVersion";

interface ProjectTableRowProps {
  id?: number;
  name: string;
  numberOfVersions?: number;
  createdAt: string;
}

export const TableRow: FC<ProjectTableRowProps> = ({
  id,
  name,
  numberOfVersions,
  createdAt,
}) => {
  const formatIsoDateToLocaleString = (isoDate: string) => {
    return new Date(isoDate).toLocaleString();
  };

  const { projectId } = useParams();
  const { mutate: handleDeleteVersion } = useDeleteVersion(projectId, id);

  const createdAtDate = formatIsoDateToLocaleString(createdAt);

  return (
    <tr>
      <td style={{ paddingLeft: "1.5rem" }}>
        <Typography level="body-xs">{name}</Typography>
      </td>

      {numberOfVersions && (
        <td style={{ paddingLeft: "0.5rem" }}>
          <Typography level="body-xs">{numberOfVersions}</Typography>
        </td>
      )}

      <td style={{ paddingLeft: "0.5rem" }}>
        <Typography level="body-xs">{createdAtDate}</Typography>
      </td>

      {!numberOfVersions ? (
        <td
          style={{
            textAlign: "end",
            padding: "0.5rem 5rem",
            verticalAlign: "center",
          }}
        >
          <Button sx={{ mb: "0.5rem" }} onClick={() => handleDeleteVersion()}>Delete <Delete /></Button>
        </td>
      ) : (
        <td
          style={{
            textAlign: "end",
            padding: "0.5rem 5rem",
            verticalAlign: "center",
          }}
        >
          <Link to={`/projects/${id}/versions`}>
            <Button sx={{ mb: "0.5rem" }} href="">See all versions</Button>
          </Link>
        </td>
      )}
    </tr>
  );
};
