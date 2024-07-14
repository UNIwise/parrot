import { Button, Typography } from "@mui/joy";
import { FC } from "react";
import { Link } from "react-router-dom";
import { ManageVersionModal } from "./Modal";

interface ProjectTableRowProps {
  projectId: number;
  versionId?: number;
  name: string;
  numberOfVersions?: number;
  createdAt: string;
}

export const TableRow: FC<ProjectTableRowProps> = ({
  projectId,
  versionId,
  name,
  numberOfVersions,
  createdAt,
}) => {
  const formatIsoDateToLocaleString = (isoDate: string) => {
    return new Date(isoDate).toLocaleString();
  };

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

      {numberOfVersions ? (
        <td
          style={{
            textAlign: "end",
            padding: "0.5rem 5rem",
            verticalAlign: "center",
          }}
        >
          <Link to={`/projects/${projectId}/versions`}>
            <Button sx={{ mb: "0.5rem", backgroundColor: '#0078ff' }} href="">See all versions</Button>
          </Link>
        </td>
      ) : (
        <td
          style={{
            textAlign: "end",
            padding: "0.5rem 5rem",
            verticalAlign: "center",
          }}
        >
          <ManageVersionModal versionId={versionId} projectId={projectId!} versionName={name} />
        </td>
      )}
    </tr >
  );
};
