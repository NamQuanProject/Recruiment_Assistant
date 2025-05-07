// import React from "react";

const Footer = () => {
  return (
    <footer className="bg-navbar text-white py-3 w-full">
      <div className="container mx-auto text-center">
        <p className="text-sm">
          © {new Date().getFullYear()} Recruitment Assistant. All rights reserved.
        </p>
        <p className="text-sm mt-1">
          Built with ❤️ by our team, Solidarity.
        </p>
      </div>
    </footer>
  );
};

export default Footer;