-- +goose Up
-- <u>
-- </u>
INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 1, 1, 1,'โครงการนำเสนอความเป็นมาและวัตถุประสงค์ของการจัดงานได้อย่างชัดเจน', 
'โครงการนำเสนอ<u>ความเป็นมาและวัตถุประสงค์</u>ของการจัดงานได้อย่าง<u>ชัดเจน</u>');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 1, 2, 2,'โครงการให้รายละเอียดรูปแบบการจัดงาน (สถานที่ ประเภทงาน จำนวนผู้เข้าร่วม กิจกรรมที่จะดำเนินการ) อย่างครบถ้วน', 
'โครงการให้รายละเอียด<u>รูปแบบการจัดงาน</u> (สถานที่ ประเภทงาน จำนวนผู้เข้าร่วม กิจกรรมที่จะดำเนินการ) อย่าง<u>ครบถ้วน</u>');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 1, 3, 3,'ในภาพรวม ข้อเสนอโครงการถือว่ามีคุณภาพดี', 
'ในภาพรวม <u>ข้อเสนอโครงการ</u>ถือว่า<u>มีคุณภาพดี</u>');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 2, 1, 4,'โครงการแสดงรายละเอียดการจัดการเส้นทางวิ่งอย่างเหมาะสม เช่น ระยะทาง เส้นทางวิ่ง', 
'โครงการแสดงรายละเอียดการ<u>จัดการเส้นทางวิ่ง</u>อย่าง<u>เหมาะสม</u> เช่น ระยะทาง เส้นทางวิ่ง');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 2, 2, 5,'โครงการแสดงให้เห็นการประสานความร่วมมือกับพื้นที่จัดงาน โดยเฉพาะองค์การปกครองส่วนท้องถิ่น', 
'โครงการแสดงให้เห็นการ<u>ประสานความร่วมมือ</u>กับพื้นที่จัดงาน โดยเฉพาะ<u>องค์การปกครองส่วนท้องถิ่น</u>');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 2, 3, 6,'โครงการมีแผนการจัดการด้านความปลอดภัยและการแพทย์ เช่น บุคลากรด้านเพื่อปฐมพยาบาล รถฉุกเฉิน', 
'โครงการมีแผนการ<u>จัดการด้านความปลอดภัยและการแพทย์</u> เช่น บุคลากรด้านเพื่อปฐมพยาบาล รถฉุกเฉิน');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 2, 4, 7,'โครงการมีแผนการจัดการสิ่งอำนวยความสะดวก (เช่น อาหาร น้ำดื่ม) อย่างเหมาะสมและพอเพียง', 
'โครงการมีแผนการจัดการ<u>สิ่งอำนวยความสะดวก</u> (เช่น อาหาร น้ำดื่ม) อย่าง<u>เหมาะสมและพอเพียง</u>');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 2, 5, 8,'โครงการกำหนดวิธีการสื่อสารกับผู้เข้าร่วมเพื่อให้แน่ใจว่าการกระจายข้อมูลข่าวสารเกี่ยวกับงานรับรู้ทั่วถึง', 
'โครงการกำหนดวิธีการสื่อสารกับผู้เข้าร่วมเพื่อให้แน่ใจว่า<u>การกระจายข้อมูลข่าวสารเกี่ยวกับงานรับรู้ทั่วถึง</u>');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 2, 6, 9,'โครงการมีการออกแบบกระบวนการประเมินผลจากการจัดงานเพื่อการพัฒนาในอนาคต', 
'โครงการมีการออกแบบกระบวน<u>การประเมินผลจากการจัดงาน</u>เพื่อการพัฒนาในอนาคต');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 3, 1, 10,'โครงการมีความสอดคล้องกับหลักการ พันธกิจและแนวคิดของ สสส. ในฐานะองค์กรสร้างเสริมสุขภาพ', 
'โครงการมีความ<u>สอดคล้องกับหลักการ พันธกิจและแนวคิดของ สสส.</u> ในฐานะองค์กรสร้างเสริมสุขภาพ');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 3, 2, 11,'โครงการมีความเกี่ยวข้องกับกลุ่มเป้าหมายการดำเนินงานของ สสส. ด้านการส่งเสริมกิจกรรมทางกาย', 
'โครงการมี<u>ความเกี่ยวข้องกับกลุ่มเป้าหมาย</u>การดำเนินงานของ สสส. <u>ด้านการส่งเสริมกิจกรรมทางกาย</u>');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 3, 3, 12,'โครงการมีความสอดคล้องกับแนวคิดของแผนส่งเสริมกิจกรรมทางกาย สสส.', 
'โครงการมี<u>ความสอดคล้องกับแนวคิดของแผนส่งเสริมกิจกรรมทางกาย</u> สสส.');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 4, 1, 13,'โลโก้หรืออัตลักษณ์ของ สสส. ได้รับการนำเสนออย่างชัดเจน (เช่น Banners, ป้ายโฆษณา หรือสื่อประเภทต่างๆ)', 
'<u>โลโก้หรืออัตลักษณ์ของ สสส. ได้รับการนำเสนอ</u>อย่างชัดเจน (เช่น Banners, ป้ายโฆษณา หรือสื่อประเภทต่างๆ)');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 4, 2, 14,'โครงการแสดงให้เห็นโอกาส/ศักยภาพของการใช้สื่อที่หลากหลายทั้ง online และ offline ที่จะเกิด ประโยชน์ต่อการสื่อสารภาพลักษณ์สสส.', 
'โครงการแสดงให้เห็นโอกาส/ศักยภาพของ<u>การใช้สื่อที่หลากหลายทั้ง online และ offline</u> ที่จะเกิด ประโยชน์ต่อการสื่อสารภาพลักษณ์สสส.');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 4, 3, 15,'ในการจัดงาน สสส. มีโอกาสที่จะสร้างกิจกรรมมีส่วนร่วมกับผู้เข้าร่วมงานวิ่ง เช่น บูธจัดงาน พื้นที่กิจกรรม เป็นต้น', 
'ในการจัดงาน สสส. <u>มีโอกาสที่จะสร้างกิจกรรมมีส่วนร่วมกับผู้เข้าร่วมงานวิ่ง</u> เช่น บูธจัดงาน พื้นที่กิจกรรม เป็นต้น');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 5, 1, 16,'ในกรณีที่มีการจัดงานมากกว่าหนึ่งครั้ง งานแข่งขันวิ่งเพื่อสุขภาพนี้มีภาพลักษณ์การจัดงานที่ดีและ ประสบความสำเร็จ (หากเป็นการจัดงานครั้งแรกให้ระบุเกณฑ์ประเมิน 1)', 
'ในกรณีที่มีการจัดงานมากกว่าหนึ่งครั้ง งานแข่งขันวิ่งเพื่อสุขภาพนี้มี<u>ภาพลักษณ์การจัดงานที่ดีและ ประสบความสำเร็จ</u> (หากเป็นการจัดงานครั้งแรกให้ระบุเกณฑ์ประเมิน 1)');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 5, 2, 17,'ผู้รับผิดชอบโครงการและผู้จัดงานมีประวัติที่ดีและมีประสบการณ์ที่น่าเชื่อถือจากการจัดงานวิ่ง เพื่อสุขภาพ ที่ผ่าน ๆ มา', 
'<u>ผู้รับผิดชอบโครงการและผู้จัดงาน</u>มีประวัติที่ดีและมีประสบการณ์ที่<u>น่าเชื่อถือ</u>จากการจัดงานวิ่ง เพื่อสุขภาพ ที่ผ่าน ๆ มา');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 6, 1, 18,'ภาพรวมโครงการเป็นประโยชน์เหมาะสมแก่ สสส. ที่จะสนับสนุนทุนอุปถัมภ์', 
'ภาพรวมโครงการเป็นประโยชน์<u>เหมาะสม</u>แก่ สสส. <u>ที่จะสนับสนุนทุนอุปถัมภ์</u>');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 6, 2, 19,'สสส. มีโอกาสที่จะสร้างความร่วมมือหรือขยายเครือข่ายร่วมกับองค์กรที่ขอรับการสนับสนุนในระยะยาว', 
'สสส. <u>มีโอกาสที่จะสร้างความร่วมมือหรือขยายเครือข่าย</u>ร่วมกับองค์กรที่ขอรับการสนับสนุนในระยะยาว');

INSERT INTO review_criteria (criteria_version, group_number, in_group_number, order_number, question_text, display_text)
VALUES (1, 6, 3, 20,'การสนับสนุนทุนอุปถัมภ์ครั้งนี้เป็นประโยชน์ต่อสสส.ในหลากหลายมิติ', 
'การสนับสนุนทุนอุปถัมภ์ครั้งนี้<u>เป็นประโยชน์ต่อสสส.</u>ในหลากหลายมิติ');

-- +goose Down

DELETE FROM review_criteria WHERE criteria_version = 1;