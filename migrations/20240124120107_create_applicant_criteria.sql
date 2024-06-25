-- +goose Up
CREATE TABLE applicant_criteria(
  id SERIAL PRIMARY KEY NOT NULL,
  code VARCHAR(64) NOT NULL,
  criteria_version SMALLINT NOT NULL,
  order_number SMALLINT NOT NULL,
  display VARCHAR(512) NOT NULL,
  pdf_display VARCHAR(512) NOT NULL
);
CREATE TABLE applicant_score (
  id SERIAL PRIMARY KEY NOT NULL,
  project_history_id INT REFERENCES project_history(id),
  applicant_criteria_id INT REFERENCES applicant_criteria(id),
  score SMALLINT NOT NULL
);
-- SEED DATA
INSERT INTO applicant_criteria (
    code,
    criteria_version,
    order_number,
    display,
    pdf_display
  )
VALUES (
    'project_self_score',
    1,
    1,
    '1. แผนการจัดระบบน้ำดื่ม คำนวณ<u>ปริมาณน้ำดื่มสะอาดไว้อย่างเพียงพอ</u>ในทุกจุดที่ให้บริการน้ำดื่มสำหรับนักวิ่งในทุกระยะทาง',
    'คำนวณปริมาณน้ำดื่มไว้อย่างเพียงพอในทุกจุดบริการ'
  );
INSERT INTO applicant_criteria (
    code,
    criteria_version,
    order_number,
    display,
    pdf_display
  )
VALUES (
    'project_self_score',
    1,
    2,
    '2. แผนการจัด<u>อาหาร</u>เพื่อบริการสำหรับนักวิ่ง จัดเตรียมตามหลักสุขอนามัยทุกประการไม่ก่อให้เกิดอาการเจ็บป่วยหลังรับประทาน',
    'จัดเตรียมอาหารตามหลักสุขอนามัยทุกประการ ไม่ก่อให้เกิดอาการเจ็บป่วยหลังรับประทาน'
  );
INSERT INTO applicant_criteria (
    code,
    criteria_version,
    order_number,
    display,
    pdf_display
  )
VALUES (
    'project_self_score',
    1,
    3,
    '3. <u>แผนการจัดการขยะ</u>จะไม่สร้างขยะตกค้างในบริเวณจัดกิจกรรมหลังกิจกรรมสิ้นสุด',
    'ไม่สร้างขยะตกค้างในบริเวณจัดกิจกรรมหลังกิจกรรมสิ้นสุด'
  );
INSERT INTO applicant_criteria (
    code,
    criteria_version,
    order_number,
    display,
    pdf_display
  )
VALUES (
    'project_self_score',
    1,
    4,
    '4. <u>แผนการจัดระบบแสงสว่างที่เพียงพอในบริเวณจุดเริ่มและเส้นชัย</u> ตลอดจนตลอดเส้นทางวิ่งทุกระยะให้นักวิ่งทุกคน สามารถมองเห็นเส้นทางโดยเฉพาะช่วงเช้ามืดหรือกลางคืน',
    'แสงสว่างเพียงพอในบริเวณจุดเริ่มต้นและเส้นชัย ตลอดจนตลอดเส้นทางวิ่งทุกระยะ'
  );
INSERT INTO applicant_criteria (
    code,
    criteria_version,
    order_number,
    display,
    pdf_display
  )
VALUES (
    'project_self_score',
    1,
    5,
    '5. แผนสุขาภิบาล มีการประสานงานให้มีบริการ<u>ห้องน้ำตามจุดที่เหมาะสมของเส้นทางวิ่งแต่ละระยะ</u>',
    'มีบริการห้องน้ำตามจุดที่เหมาะสมของเส้นทางวิ่งแต่ละระยะ'
  );
INSERT INTO applicant_criteria (
    code,
    criteria_version,
    order_number,
    display,
    pdf_display
  )
VALUES (
    'project_self_score',
    1,
    6,
    '6. แผนการลงทะเบียนรายงานตัวมี<u>การจัดการระบบการรับลงทะเบียนเพื่อรับเบอร์วิ่ง (BIB) และ/หรือ Chip Time เสื้อที่ระลึก และอุปกรณ์อื่น ๆ</u> อย่างเหมาะสม และมีจำนวนเพียงพอ สำหรับนักวิ่งทุกคน',
    'จัดการระบบรับลงทะเบียน เพื่อรับ BIB, Chip Time, เสื้อที่ระลึก และอุปกรณ์อื่น ๆ อย่างเหมาะสม และมีจำนวนเพียงพอสำหรับนักวิ่งทุกคน'
  );
INSERT INTO applicant_criteria (
    code,
    criteria_version,
    order_number,
    display,
    pdf_display
  )
VALUES (
    'project_self_score',
    1,
    7,
    '7. <u>การรับฝากของส่วนตัวของนักวิ่ง</u> มีเจ้าหน้าที่ สถานที่ และอุปกรณ์พร้อมสำหรับการรับฝากของ',
    'มีเจ้าหน้าที่ สถานที่ และอุปกรณ์พร้อมสำหรับการรับฝากของ'
  );
INSERT INTO applicant_criteria (
  code,
  criteria_version,
  order_number,
  display,
  pdf_display
)
VALUES (
    'project_self_score',
    1,
    8,
    '8. <u>แผนความปลอดภัยด้านสุขภาพ</u>ของนักวิ่ง พร้อมปฏิบัติได้จริง เช่น ตั้งจุดปฐมพยาบาลพร้อมเวชภัณฑ์ และมีการประสานงานขอทีมแพทย์/พยาบาลเคลื่อนที่',
    'แผนความปลอดภัยด้านสุขภาพของนักวิ่ง พร้อมปฏิบัติได้จริง'
  );
INSERT INTO applicant_criteria (
    code,
    criteria_version,
    order_number,
    display,
    pdf_display
  )
VALUES (
    'project_self_score',
    1,
    9,
    '9. แผนการจัดที่พัก ได้เสนอ/ประสาน ให้เกิดทางเลือกเรื่อง<u>ที่พัก</u>ให้กับนักวิ่ง และผู้ติดตามอย่างเพียงพอ (ถ้าไม่มีให้ระบุ “1”)',
    'ได้ประสานให้เกิดทางเลือกเรื่องที่พักให้กับนักวิ่งและผู้ติดตามอย่างเพียงพอ'
  );
INSERT INTO applicant_criteria (
    code,
    criteria_version,
    order_number,
    display,
    pdf_display
  )
VALUES (
    'project_self_score',
    1,
    10,
    '10. แผนการจราจร การจัดเตรียม/ประสานให้มี<u>พื้นที่จอดรถ</u>ทั้งรถจักรยานยนต์และรถยนต์ให้กับนักวิ่ง',
    'จัดเตรียมให้มีพื้นที่จอดรถทั้งรถจักรยานยนต์และรถยนต์ให้กับนักวิ่ง'
  );
INSERT INTO applicant_criteria (
    code,
    criteria_version,
    order_number,
    display,
    pdf_display
  )
VALUES (
    'project_self_score',
    1,
    11,
    '11. แผนเสนอข้อมูลเส้นทางวิ่ง <u>โดยจัดทำแผนที่/แผนผังเส้นทางวิ่งทุกระยะ</u> ในบริเวณพื้นที่จัดกิจกรรมเพื่อสื่อสารได้อย่างชัดเจน',
    'จัดทำแผนผังเส้นทางวิ่งทุกระยะในบริเวณพื้นที่จัดกิจกรรม'
  );
-- +goose Down
ALTER TABLE applicant_score DROP COLUMN project_history_id;
ALTER TABLE applicant_score DROP COLUMN applicant_criteria_id;
DROP TABLE applicant_score;
DROP TABLE applicant_criteria;