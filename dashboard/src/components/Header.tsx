import styled from "@emotion/styled";
import { Button, Sheet, Stack, Typography } from "@mui/joy";
import { FC, useEffect, useRef, useState } from "react";
import { Link } from "react-router-dom";
import soundEffect from "../assets/slap.mp3";
import { ColorSchemeToggle } from "./ColorSchemeToggle";

const rainbowGradient = `linear-gradient(
  to right,
  rgba(255, 0, 0, 0.7),
  rgba(255, 165, 0, 0.7),
  rgba(255, 255, 0, 0.7),
  rgba(0, 128, 0, 0.7),
  rgba(0, 0, 255, 0.7),
  rgba(128, 0, 128, 0.7)
)`;

const AnimatedSheet = styled(Sheet, {
  shouldForwardProp: (prop) => prop !== "isblue" && prop !== "isanimating",
})<{
  isblue: boolean;
  isanimating: boolean;
}>`
  transition: background 0.5s linear;
  background-size: 100% 100%;
  background-position: ${(props) =>
    props.isblue ? "left bottom" : "right bottom"};
  background-color: ${(props) =>
    props.isanimating
      ? `linear-gradient(to right, blue 50%, ${rainbowGradient} 50%)`
      : props.isblue
        ? "#4393E4"
        : rainbowGradient};
  background-image: ${(props) =>
    props.isanimating
      ? `linear-gradient(to right, blue 50%, ${rainbowGradient} 50%)`
      : props.isblue
        ? "#4393E4"
        : rainbowGradient};
`;

type HeaderProps = {
  items?: {
    name: string;
    to?: string;
  }[];
};

export const Header: FC<HeaderProps> = ({ items }) => {
  const [isBlue, setIsBlue] = useState(true);
  const [isAnimating, setIsAnimating] = useState(false);
  const audioRef = useRef<HTMLAudioElement | null>(null);

  const toggleBackground = () => {
    setIsAnimating(true);
    if (audioRef.current) {
      audioRef.current.play();
    }
  };

  useEffect(() => {
    if (isAnimating) {
      const timer = setTimeout(() => {
        setIsBlue((prev) => !prev);
        setIsAnimating(false);
      }, 1000);
      return () => clearTimeout(timer);
    }
  }, [isAnimating]);

  return (
    <AnimatedSheet
      isblue={isBlue}
      isanimating={isAnimating}
      sx={{
        p: 2,
        mb: 4,
        display: "flex",
        justifyContent: "space-between",
        borderRadius: "5px",
        alignItems: "center",
        height: 64,
      }}
    >
      <Stack direction="row" spacing={2}>
        {items?.map((item, index) => (
          <>
            {index > 0 && (
              <Typography
                level="h2"
                sx={{
                  color: "primary.solidColor",
                  fontSize: "1.75em",
                  textDecoration: "none",
                }}
              >
                /
              </Typography>
            )}

            <Typography
              key={index}
              component={Link}
              to={item.to || "/"}
              level="h2"
              sx={{
                color: "primary.solidColor",
                fontSize: "1.75em",
                textDecoration: "none",
                "&:hover": {
                  textDecoration: "underline",
                },
              }}
            >
              {item.name}
            </Typography>
          </>
        ))}
      </Stack>

      <div>
        <Button
          onClick={toggleBackground}
          sx={{ backgroundColor: "transparent", height: 20 }}
          disabled={isAnimating}
        />

        <ColorSchemeToggle sx={{ ml: "2rem" }} />
      </div>
      <audio ref={audioRef} src={soundEffect} />
    </AnimatedSheet>
  );
};
