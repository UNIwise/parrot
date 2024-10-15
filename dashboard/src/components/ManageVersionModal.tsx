import { Delete } from "@mui/icons-material";
import Add from "@mui/icons-material/Add";
import { FormLabel } from "@mui/joy";
import Button from "@mui/joy/Button";
import DialogContent from "@mui/joy/DialogContent";
import DialogTitle from "@mui/joy/DialogTitle";
import FormControl from "@mui/joy/FormControl";
import Input from "@mui/joy/Input";
import Modal from "@mui/joy/Modal";
import ModalDialog from "@mui/joy/ModalDialog";
import Stack from "@mui/joy/Stack";
import { ChangeEvent, FC, FormEvent, useState } from "react";
import { useDeleteVersion } from "../api/hooks/useDeleteVersion";
import { usePostVersion } from "../api/hooks/usePostVersion";

type ManageVersionModalProps = {
  projectId: number;
  versionId?: number;
  versionName?: string;
};

export const ManageVersionModal: FC<ManageVersionModalProps> = ({
  projectId,
  versionId,
  versionName,
}) => {
  const [open, setOpen] = useState(false);
  const [newProjectVersion, setNewProjectVersion] = useState("");
  const [validProjectVersion, setValidProjectVersion] = useState(false);

  const { mutate: deleteVersion } = useDeleteVersion(projectId, versionId!);
  const { mutate: postNewVersion } = usePostVersion(projectId);

  const onProjectNameChange = (event: ChangeEvent<HTMLInputElement>) => {
    const valid = /^[a-zA-Z0-9.-_]*$/.test(event.target.value);
    setValidProjectVersion(valid);

    setNewProjectVersion(event.target.value);
  };

  const onSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (versionId) {
      deleteVersion();
    } else {
      postNewVersion({ name: newProjectVersion });
    }

    setOpen(false);
  };

  return (
    <>
      {versionId ? (
        <Button
          variant="outlined"
          color="danger"
          sx={{ color: (t) => t.palette.danger[500] }}
          endDecorator={<Delete />}
          onClick={() => setOpen(true)}
        >
          Delete
        </Button>
      ) : (
        <Button
          color="primary"
          startDecorator={<Add />}
          onClick={() => setOpen(true)}
          sx={{
            height: "36px",
          }}
        >
          New Version
        </Button>
      )}

      <Modal open={open} onClose={() => setOpen(false)}>
        <ModalDialog>
          <form onSubmit={onSubmit}>
            {versionId ? (
              <Stack spacing={2}>
                <DialogTitle>Delete version</DialogTitle>
                <DialogContent>
                  Are you sure you want to delete {versionName} version?
                </DialogContent>

                <DialogContent>
                  Deletion of the version will be done in the background and
                  will take a couple of seconds.
                </DialogContent>

                <Button type="submit" color="danger">
                  Delete
                </Button>
              </Stack>
            ) : (
              <Stack spacing={2}>
                <DialogTitle>New version</DialogTitle>

                <DialogContent>
                  Creating the version will be done in the background and will
                  take a couple of minutes.
                </DialogContent>

                <FormControl>
                  <FormLabel>Version tag</FormLabel>
                  <Input
                    autoFocus
                    required
                    onChange={onProjectNameChange}
                    value={newProjectVersion}
                    error={!!newProjectVersion && !validProjectVersion}
                  />
                </FormControl>

                <Button type="submit" disabled={!validProjectVersion}>
                  Add version
                </Button>
              </Stack>
            )}
          </form>
        </ModalDialog>
      </Modal>
    </>
  );
};
