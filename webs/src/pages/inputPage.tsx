import NavBar from "../components/navbar";
import IPBox from "../components/inputbox";
import Criteria from "../components/criteria";

import Logo1 from '../assets/Ellipse 1.png';
import Logo2 from '../assets/Ellipse 2.png';
import Logo3 from '../assets/Ellipse 3.png';
import Logo4 from '../assets/Ellipse 4.png';

const inputPage = () => {
  return (
    <>
    <NavBar />
    <IPBox />
    <Criteria />
    <img src={Logo1} alt="Logo 1" className="absolute bottom-0 right-0 translate-y-2/5 -z-10" />
    <img src={Logo2} alt="Logo 2" className="absolute top-0 right-0 -translate-x-[200px] -z-10" />
    <img src={Logo3} alt="Logo 3" className="absolute bottom-0 left-0 translate-y-[400px] translate-x-[30px] -z-10" />
    <img src={Logo4} alt="Logo 4" className="absolute top-0 left-0 translate-y-[30px] -z-10" /> 
    </>
  );
}

export default inputPage;