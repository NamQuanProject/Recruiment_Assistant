import NavBar from "../components/navbar";
import IPBox from "../components/inputbox";
import Criteria from "../components/criteria";
import { useState } from "react";
import Logo1 from '../assets/Ellipse 1.png';
import Logo2 from '../assets/Ellipse 2.png';
import Logo3 from '../assets/Ellipse 3.png';
import Logo4 from '../assets/Ellipse 4.png';
import Footer from "../components/footer";
const InputPage = () => {
  // const [criteriaData, setCriteriaData] = useState<any>(null); // State to store the response data
  const [loading, setLoading] = useState(false); // State to manage loading

  return (
    <>
      <NavBar />
      {/* Pass setLoading and loading to IPBox */}
      <IPBox setLoading={setLoading} />
      {/* Pass loading to Criteria */}
      <Criteria  loading={loading} />
      <img src={Logo1} alt="Logo 1" className="fixed bottom-0 right-0 translate-y-2/5 -z-10" />
      <img src={Logo2} alt="Logo 2" className="fixed top-0 right-0 -translate-x-[200px] -z-10" />
      <img src={Logo3} alt="Logo 3" className="fixed bottom-0 left-0 translate-x-[30px] -z-10" />
      <img src={Logo4} alt="Logo 4" className="fixed top-0 left-0 translate-y-[30px] -z-10" />
      <Footer />
    </>
  );
};

export default InputPage;