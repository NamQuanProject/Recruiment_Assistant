import Logo from "../assets/react.svg";
const NavBar = () => {
  const flexBetween = "flex justify-between items-center";
  return (
    <div className={`${flexBetween} fixed top-0 z-30 w-full py-6 border-2 border-gray-400 shadow-md bg-cyan-50`}>
      <div className={`${flexBetween} mx-auto w-5/6`}>
        <div className={`${flexBetween} w-full gap-12`} >
          <img alt="logo" src={Logo}/>
          <div className={`${flexBetween} w-full`}>
              <div className={`${flexBetween} gap-8 text-sm`}>
                <div className="cursor-pointer">Home</div>
                <div className="cursor-pointer">Products</div>
                <div className="cursor-pointer">Contact</div>
              </div>
              <div className={`${flexBetween} gap-8`}>
                <div className="cursor-pointer">Employer Login</div>
                <div className="cursor-pointer">Employee Login</div>
              </div>  
          </div>
        </div>
      </div>
    </div>
  );
};

export default NavBar;