-- +goose Up
CREATE TABLE applicant_criteria(
  id SERIAL PRIMARY KEY NOT NULL,
  code VARCHAR(64) NOT NULL,
  criteria_version SMALLINT NOT NULL,
  order_number SMALLINT NOT NULL,
  display VARCHAR(512) NOT NULL
);
INSERT INTO applicant_criteria (code, criteria_version, order_number, display)
VALUES (
    'project_self_score',
    1,
    1,
    '1. งานวิ่งที่ท่านกำลังจะจัดได้<u>คำนวณปริมาณน้ำดื่มสะอาดไว้อย่างเพียงพอ</u>สำหรับนักวิ่งทุกคนในทุกระยะตลอดการแข่งขัน'
  );
INSERT INTO applicant_criteria (code, criteria_version, order_number, display)
VALUES (
    'project_self_score',
    1,
    2,
    '2. <u>อาหาร</u>ที่จัดบริการสำหรับนักวิ่งจัดเตรียมอย่าง<u>สะอาด</u> ไม่ก่อให้เกิดอาการเจ็บป่วยหลังรับประทาน'
  );
INSERT INTO applicant_criteria (code, criteria_version, order_number, display)
VALUES (
    'project_self_score',
    1,
    3,
    '3. <u>แผนการจัดการขยะ</u>จะไม่สร้างขยะตกค้างในบริเวณจัดงานหลังกิจกรรมสิ้นสุด'
  );
INSERT INTO applicant_criteria (code, criteria_version, order_number, display)
VALUES (
    'project_self_score',
    1,
    4,
    '4. งานวิ่งที่ท่านกำลังจะจัดมี<u>แสงสว่างที่เพียงพอ</u>แก่การมองเห็นเส้นทางสำหรับนักวิ่งทุกคน โดยเฉพาะช่วงเช้ามืดหรือกลางคืน'
  );
INSERT INTO applicant_criteria (code, criteria_version, order_number, display)
VALUES (
    'project_self_score',
    1,
    5,
    '5. งานวิ่งที่ท่านกำลังจะจัดมีแผนหรือประสานงานให้มีบริการ<u>ห้องน้ำตลอดเส้นทางวิ่ง</u>'
  );
INSERT INTO applicant_criteria (code, criteria_version, order_number, display)
VALUES (
    'project_self_score',
    1,
    6,
    '6. <u>การรับของที่ระลึก เช่น เสื้อ เหรียญรางวัล</u> มีจุดบริการและเจ้าหน้าที่รองรับอย่างเหมาะสม และมีจำนวนที่จัดเตรียมเพียงพอสำหรับนักวิ่งทุกคน'
  );
INSERT INTO applicant_criteria (code, criteria_version, order_number, display)
VALUES (
    'project_self_score',
    1,
    7,
    '7. งานวิ่งที่ท่านกำลังจะจัดมี<u>แผนรองรับเรื่องความปลอดภัย</u>ของนักวิ่งที่พร้อมปฏิบัติได้จริง'
  );
INSERT INTO applicant_criteria (code, criteria_version, order_number, display)
VALUES (
    'project_self_score',
    1,
    8,
    '8. งานวิ่งที่ท่านกำลังจะจัดได้เสนอ/ประสานให้เกิดทางเลือกเรื่อง<u>ที่พัก</u>ให้กับผู้เข้าร่วมร่วมกิจกรรมอย่างเพียงพอ (ถ้าไม่มีให้ระบุ “1”)'
  );
INSERT INTO applicant_criteria (code, criteria_version, order_number, display)
VALUES (
    'project_self_score',
    1,
    9,
    '9. งานวิ่งที่ท่านกำลังจะจัดมีการจัดเตรียม/ประสานให้มี<u>พื้นที่จอดรถ</u>ทั้งรถจักรยานยนต์และรถยนต์ให้กับผู้เข้าร่วมกิจกรรม'
  );
INSERT INTO applicant_criteria (code, criteria_version, order_number, display)
VALUES (
    'project_self_score',
    1,
    10,
    '10. งานวิ่งที่ท่านกำลังจะจัดแสดง<u>แผนที่/แผนผังเส้นทางกิจกรรม</u>ในบริเวณพื้นที่จัดงานเพื่อสื่อสารได้อย่างชัดเจน'
  );

-- +goose Down
DROP TABLE applicant_criteria;