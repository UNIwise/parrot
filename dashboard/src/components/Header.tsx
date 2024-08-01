import styled from "@emotion/styled";
import { Button, Sheet, Typography } from "@mui/joy";
import React, { useEffect, useRef, useState } from "react";
import { Link } from "react-router-dom";
import { ColorSchemeToggle } from "./ColorSchemeToggle";

import soundEffect from '../assets/slap.mp3';

const rainbowGradient = `linear-gradient(
  to right,
  rgba(255, 0, 0, 0.7),
  rgba(255, 165, 0, 0.7),
  rgba(255, 255, 0, 0.7),
  rgba(0, 128, 0, 0.7),
  rgba(0, 0, 255, 0.7),
  rgba(128, 0, 128, 0.7)
)`;

const AnimatedSheet = styled(Sheet, { label: 'styledComponent' }) <{ isblue: boolean, isanimating: boolean }>`
  transition: background 0.5s linear;
  background-size: 100% 100%;
  background-position: ${props => props.isblue ? "left bottom" : "right bottom"};
  background-color: ${props =>
    props.isanimating
      ? `linear-gradient(to right, blue 50%, ${rainbowGradient} 50%)`
      : props.isblue
        ? "#4393E4"
        : rainbowGradient
  };
    background-image: ${props =>
    props.isanimating
      ? `linear-gradient(to right, blue 50%, ${rainbowGradient} 50%)`
      : props.isblue
        ? "#4393E4"
        : rainbowGradient
  };
`

export const Header: React.FC = () => {
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
        setIsBlue(prev => !prev);
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
        mb: 6,
        display: "flex",
        justifyContent: "space-between",
        borderRadius: "5px",
        alignItems: "center",
        height: 64,
      }}
    >
      <Typography
        component={Link}
        to={"/"}
        level="h4"
        sx={{
          fontWeight: "bold",
          color: "primary.solidColor",
          fontSize: "2.2em",
        }}
      >
        Parrot
      </Typography>

      <div>
        <Button
          onClick={toggleBackground}
          sx={{ backgroundColor: 'transparent', height: 20 }}
          disabled={isAnimating}

        />
        <ColorSchemeToggle sx={{ ml: "2rem" }} />
      </div>
      <audio ref={audioRef} src={soundEffect} />
    </AnimatedSheet>
  );
};
