import Logo from "../assets/logo.png";
import { Link } from "react-router-dom";
const NavBar = () => {
  const flexBetween = "flex justify-between items-center";
  return (
    <div className={`${flexBetween} fixed top-0 z-30 w-full py-5 border-2 border-gray-400 shadow-md bg-navbar`}>
      <div className={`${flexBetween} mx-auto w-5/6`}>
        <div className={`${flexBetween} w-full gap-12`} >
          <img alt="logo" className="object-contain w-24 h-12 " src={Logo}/>
          <div className={`${flexBetween} w-full`}>
              <div className={`${flexBetween} gap-8 text-sm`}>
              <Link to="/" className="cursor-pointer text-base text-gray-200">Home</Link> {/* Navigate to / */}
                <div className="cursor-pointer text-base text-gray-200">Products</div>
                <div className="cursor-pointer text-base text-gray-200">Contact</div>
              </div>
              <div className={`${flexBetween} gap-8`}>
                <div className="cursor-pointer text-base text-gray-200">Employer</div>
                <div className="cursor-pointer text-base text-gray-200">Employee</div>
              </div>  
          </div>
        </div>
      </div>
    </div>
  );
};

export default NavBar;